package data

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	gamingAPI "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/room/internal/biz"
	etcdAPI "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"math/rand"
	"net/url"
	"time"
)

type roomRepo struct {
	data *Data
	log  *log.Helper
}

type Room struct {
	Id           int64     `gorm:"primary_key;auto_increment"`
	Name         string    `gorm:"size:255"`
	Host         int64     `gorm:"not null"`
	Private      bool      `gorm:"not null"`
	PasswordHash string    `gorm:"size:255"`
	PasswordReal string    `gorm:"size:16"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

type RoomMember struct {
	Id       int64     `gorm:"primary_key;auto_increment"`
	RoomId   int64     `gorm:"not null;"`
	UserId   int64     `gorm:"not null;"`
	Role     int       `gorm:"not null"`
	JoinedAt time.Time `gorm:"column:joined_at"`
}

type JoinedRoom struct {
	Name         string    `gorm:"size:255"`
	Host         int64     `gorm:"not null"`
	Private      bool      `gorm:"not null"`
	PasswordHash string    `gorm:"size:255"`
	PasswordReal string    `gorm:"size:16"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	RoomId       int64     `gorm:"not null;"`
	UserId       int64     `gorm:"not null;"`
	Role         int       `gorm:"not null"`
	JoinedAt     time.Time `gorm:"column:joined_at"`
}

func NewRoomRepo(data *Data, logger log.Logger) biz.RoomRepo {
	return &roomRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/room")),
	}
}

func (r *roomRepo) CreateRoom(ctx context.Context, room *biz.Room) error {
	roomModel := Room{
		Name:    room.Name,
		Host:    room.Host,
		Private: room.Private,
	}
	if roomModel.Private {
		hashPassword := hex.EncodeToString(md5.New().Sum([]byte(room.Password)))
		roomModel.PasswordHash = hashPassword
		roomModel.PasswordReal = room.Password
	}
	return r.data.db.Model(&roomModel).WithContext(ctx).Create(&roomModel).Error
}

func (r *roomRepo) GetRoom(ctx context.Context, id int64) (*biz.Room, error) {
	room := Room{}
	if err := r.data.db.Model(&room).WithContext(ctx).Where("id =?", id).First(&room).Error; err != nil {
		return nil, err
	}
	var memberCount int64 = 0
	err := r.data.db.Model(&Room{}).WithContext(ctx).Where("room_id=?", id).Count(&memberCount).Error
	if err != nil {
		return nil, err
	}
	return room.ToBizRoom(), nil
}

func (r *roomRepo) ListRooms(ctx context.Context, page int, pageSize int) ([]*biz.Room, int, error) {
	var rooms []*Room
	var total int64 = 0
	err := r.data.db.Model(&Room{}).Count(&total).WithContext(ctx).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.data.db.Model(&Room{}).
		WithContext(ctx).
		Offset(page * pageSize).
		Limit(pageSize).
		Find(&rooms).
		Error
	if err != nil {
		return nil, 0, err
	}
	result := make([]*biz.Room, 0, len(rooms))
	for _, room := range rooms {
		result = append(result, room.ToBizRoom())
	}
	return result, int(total), nil
}

func (r *roomRepo) ListJoinedRooms(ctx context.Context, userId int64, page int, pageSize int) ([]*biz.JoinedRoom, int, error) {
	var joinedRooms []*JoinedRoom
	var total int64 = 0
	err := r.data.db.Model(&RoomMember{}).Where("user_id =?", userId).Count(&total).WithContext(ctx).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.data.db.Model(&RoomMember{}).
		InnerJoins("room on room.id = room_member.room_id").
		Where("room_member.user_id = ?", userId).
		Offset(page * pageSize).
		Limit(pageSize).
		Find(&joinedRooms).
		WithContext(ctx).Error
	if err != nil {
		return nil, 0, err
	}
	result := make([]*biz.JoinedRoom, 0, len(joinedRooms))
	for _, joinedRoom := range joinedRooms {
		result = append(result, joinedRoom.ToBizJoinedRoom())
	}
	return result, int(total), nil
}

func (r *roomRepo) GetRoomMember(ctx context.Context, roomId int64, userId int64) (*biz.RoomMember, error) {
	member := RoomMember{}
	err := r.data.db.Model(&member).
		Where("room_id = ? AND user_id = ?", roomId, userId).
		WithContext(ctx).
		First(&member).
		Error
	if err != nil {
		return nil, err
	}
	return member.ToBizRoomMember(), nil
}

func roomSessionKey(roomId int64) string {
	return fmt.Sprintf("nesgo/room/session/%d", roomId)
}

func roomSessionLockKey(roomId int64) string {
	return fmt.Sprintf("nesgo/room/session/lock/%d", roomId)
}

// GetOrCreateRoomSession 获取或创建房间会话
// 从所有模拟器节点选择一个作为目标节点，在目标节点创建模拟器实例，并保存房间ID与目标节点映射
func (r *roomRepo) GetOrCreateRoomSession(ctx context.Context, roomId int64) (*biz.RoomSession, bool, error) {
	key := roomSessionKey(roomId)
	lockKey := roomSessionLockKey(roomId)
	// 分布式锁，避免同时创建多个session
	lockSession, err := concurrency.NewSession(r.data.etcdCli)
	if err != nil {
		return nil, false, err
	}
	defer lockSession.Close()
	locker := concurrency.NewLocker(lockSession, lockKey)
	locker.Lock()
	defer locker.Unlock()

	resp, err := r.data.etcdCli.KV.Get(ctx, key)
	if err != nil {
		return nil, false, err
	}
	if resp.Count > 0 {
		var session biz.RoomSession
		_ = json.Unmarshal(resp.Kvs[0].Value, &session)
		return &session, true, nil
	}
	session := &biz.RoomSession{
		RoomId: roomId,
	}
	// 获取所有模拟器服务节点，随机选择
	serviceNodes, err := r.data.discovery.GetService(ctx, "nesgo.service.gaming")
	if err != nil {
		return nil, false, err
	}
	// TODO 更好的选择策略
	session.Endpoint = serviceNodes[rand.Intn(len(serviceNodes))].Endpoints[0]
	u, _ := url.Parse(session.Endpoint)
	session.Endpoint = u.Host
	// 在目标服务器创建模拟器实例
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(session.Endpoint))
	if err != nil {
		return nil, false, err
	}
	defer conn.Close()
	instance, err := gamingAPI.NewGamingClient(conn).CreateGameInstance(ctx, &gamingAPI.CreateGameInstanceRequest{
		RoomId: roomId,
		Game:   "SuperMario.nes",
	})
	if err != nil {
		return nil, false, err
	}
	// 保存session，使用目标节点返回的lease保活
	bytes, _ := json.Marshal(session)
	_, err = r.data.etcdCli.KV.Put(ctx, key, string(bytes), etcdAPI.WithLease(etcdAPI.LeaseID(instance.LeaseId)))
	return session, false, nil
}

// GetRoomSession 获取房间会话，返回模拟器节点地址
func (r *roomRepo) GetRoomSession(ctx context.Context, roomId int64) (*biz.RoomSession, error) {
	key := roomSessionKey(roomId)
	lockKey := roomSessionLockKey(roomId)
	lockSession, err := concurrency.NewSession(r.data.etcdCli)
	if err != nil {
		return nil, err
	}
	defer lockSession.Close()
	locker := concurrency.NewLocker(lockSession, lockKey)
	locker.Lock()
	defer locker.Unlock()

	resp, err := r.data.etcdCli.KV.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if resp.Count > 0 {
		var session biz.RoomSession
		_ = json.Unmarshal(resp.Kvs[0].Value, &session)
		return &session, nil
	}
	return nil, nil
}

func (r *roomRepo) RemoveRoomSession(ctx context.Context, roomId int64) error {
	panic("implement me")
}

func (r *roomRepo) AddRoomMember(ctx context.Context, member *biz.RoomMember) error {
	memberModel := &RoomMember{
		Id:       0,
		RoomId:   member.RoomId,
		UserId:   member.UserId,
		Role:     member.Role,
		JoinedAt: member.JoinedAt,
	}
	return r.data.db.Model(memberModel).Create(&memberModel).WithContext(ctx).Error
}

func (m *RoomMember) ToBizRoomMember() *biz.RoomMember {
	return &biz.RoomMember{
		UserId:   m.UserId,
		RoomId:   m.RoomId,
		Role:     m.Role,
		Id:       m.Id,
		JoinedAt: m.JoinedAt,
	}
}

func (r *Room) ToBizRoom() *biz.Room {
	return &biz.Room{
		Id:       r.Id,
		Name:     r.Name,
		Host:     r.Host,
		Private:  r.Private,
		Password: r.PasswordReal,
	}
}

func (jr *JoinedRoom) ToBizJoinedRoom() *biz.JoinedRoom {
	return &biz.JoinedRoom{
		Room: biz.Room{
			Name:         jr.Name,
			Host:         jr.Host,
			Private:      jr.Private,
			Password:     jr.PasswordReal,
			PasswordHash: jr.PasswordHash,
		},
		UserId: jr.UserId,
		Role:   jr.Role,
	}
}
