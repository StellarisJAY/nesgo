package data

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	gamingAPI "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	v1 "github.com/stellarisJAY/nesgo/backend/api/room/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/room/internal/biz"
	"github.com/stellarisJAY/nesgo/backend/pkg/cache"
	etcdAPI "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"gorm.io/gorm"
)

type roomRepo struct {
	data *Data
	log  *log.Helper
}

type Room struct {
	Id           int64     `gorm:"primary_key"`
	Name         string    `gorm:"size:255"`
	Host         int64     `gorm:"not null"`
	Private      bool      `gorm:"not null"`
	PasswordHash string    `gorm:"size:255"`
	PasswordReal string    `gorm:"size:16"`
	MemberLimit  int       `gorm:"not null"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	EmulatorName string    `gorm:"not null"`
}

type RoomMember struct {
	Id       int64     `gorm:"primary_key;auto_increment"`
	RoomId   int64     `gorm:"not null;"`
	UserId   int64     `gorm:"not null;"`
	Role     int       `gorm:"not null"`
	JoinedAt time.Time `gorm:"column:joined_at"`
}

type JoinedRoom struct {
	Room
	RoomMember
}

func NewRoomRepo(data *Data, logger log.Logger) biz.RoomRepo {
	return &roomRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/room")),
	}
}

func (r *roomRepo) CreateRoom(ctx context.Context, room *biz.Room) error {
	room.Id = r.data.snowflake.Generate().Int64()
	roomModel := Room{
		Name:         room.Name,
		Host:         room.Host,
		Private:      room.Private,
		Id:           room.Id,
		MemberLimit:  room.MemberLimit,
		CreatedAt:    time.Now(),
		EmulatorName: room.EmulatorType,
	}
	if roomModel.Private {
		hashPassword := hex.EncodeToString(md5.New().Sum([]byte(room.Password)))
		roomModel.PasswordHash = hashPassword
		roomModel.PasswordReal = room.Password
	}
	member := &RoomMember{
		Id:       r.data.snowflake.Generate().Int64(),
		RoomId:   room.Id,
		UserId:   room.Host,
		Role:     int(v1.RoomRole_Host),
		JoinedAt: time.Now(),
	}
	return r.data.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&roomModel).WithContext(ctx).Create(&roomModel).Error
		if err != nil {
			return err
		}
		return tx.Model(member).Create(member).WithContext(ctx).Error
	})
}

func (r *roomRepo) GetRoom(ctx context.Context, id int64) (*biz.Room, error) {
	room, err := cache.Get[Room](ctx, r.data.rdb, roomCacheKey(id))
	if err != nil {
		if err := r.data.db.Model(&room).WithContext(ctx).Where("id =?", id).First(&room).Error; err != nil {
			return nil, err
		}
		err := cache.Set(ctx, r.data.rdb, roomCacheKey(id), room)
		if err != nil {
			r.log.Errorf("cache set room error: %v", err)
		}
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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, int(total), nil
	}
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
		InnerJoins("JOIN `rooms` on rooms.id = room_members.room_id").
		Where("room_members.user_id = ?", userId).
		Offset(page * pageSize).
		Limit(pageSize).
		Select("room_members.*, rooms.*").
		Find(&joinedRooms).
		WithContext(ctx).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, int(total), nil
	}
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

func (r *roomRepo) CountMember(ctx context.Context, roomId int64) (int64, error) {
	var count int64 = 0
	err := r.data.db.Model(&RoomMember{}).
		Where("room_id =?", roomId).
		Count(&count).
		WithContext(ctx).
		Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func roomSessionKey(roomId int64) string {
	return fmt.Sprintf("nesgo/room/session/%d", roomId)
}

func roomSessionLockKey(roomId int64) string {
	return fmt.Sprintf("nesgo/room/session/lock/%d", roomId)
}

func roomCacheKey(roomId int64) string { return fmt.Sprintf("nesgo/room/%d", roomId) }

// GetOrCreateRoomSession 获取或创建房间会话
// 从所有模拟器节点选择一个作为目标节点，在目标节点创建模拟器实例，并保存房间ID与目标节点映射
func (r *roomRepo) GetOrCreateRoomSession(ctx context.Context, roomId int64, game string, emulatorType string) (*biz.RoomSession, bool, error) {
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
	if len(serviceNodes) == 0 {
		return nil, false, fmt.Errorf("no available gaming service endpoint found")
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
		RoomId:       roomId,
		Game:         game,
		EmulatorType: emulatorType,
	})
	if err != nil {
		return nil, false, err
	}
	session.InstanceId = instance.InstanceId
	// 保存session，使用目标节点返回的lease保活
	bytes, _ := json.Marshal(session)
	_, err = r.data.etcdCli.KV.Put(ctx, key, string(bytes), etcdAPI.WithLease(etcdAPI.LeaseID(instance.LeaseId)))
	return session, false, err
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
		session.RoomId = roomId
		_ = json.Unmarshal(resp.Kvs[0].Value, &session)
		return &session, nil
	}
	return nil, nil
}

func (r *roomRepo) DeleteRoomSession(ctx context.Context, roomId int64, instanceId string) error {
	key := roomSessionKey(roomId)
	lockKey := roomSessionLockKey(roomId)
	lockSession, err := concurrency.NewSession(r.data.etcdCli)
	if err != nil {
		return err
	}
	defer lockSession.Close()
	locker := concurrency.NewLocker(lockSession, lockKey)
	locker.Lock()
	defer locker.Unlock()

	resp, err := r.data.etcdCli.KV.Get(ctx, key)
	if err != nil {
		return err
	}
	if resp.Count == 0 {
		return nil
	}
	var session biz.RoomSession
	_ = json.Unmarshal(resp.Kvs[0].Value, &session)
	// 避免错误删除其他会话
	if instanceId != "" && session.InstanceId != instanceId {
		return nil
	}

	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(session.Endpoint))
	if err != nil {
		return err
	}
	defer conn.Close()
	gamingCli := gamingAPI.NewGamingClient(conn)
	_, err = gamingCli.DeleteGameInstance(ctx, &gamingAPI.DeleteGameInstanceRequest{
		RoomId: roomId,
		Force:  false,
	})
	if err != nil {
		return err
	}
	_, err = r.data.etcdCli.KV.Delete(ctx, key)
	return err
}

func (r *roomRepo) GetOwnedRoom(ctx context.Context, name string, host int64) (*biz.Room, error) {
	roomModel := &Room{}
	err := r.data.db.Model(roomModel).Where("name = ? AND host = ?", name, host).WithContext(ctx).First(&roomModel).Error
	if err != nil {
		return nil, err
	}
	return roomModel.ToBizRoom(), nil
}

func (r *roomRepo) ListMembers(ctx context.Context, roomId int64) ([]*biz.RoomMember, error) {
	var members []*RoomMember
	err := r.data.db.Model(&RoomMember{}).
		Where("room_id = ?", roomId).
		WithContext(ctx).
		Find(&members).
		Error
	if err != nil {
		return nil, err
	}
	result := make([]*biz.RoomMember, 0, len(members))
	for _, member := range members {
		result = append(result, member.ToBizRoomMember())
	}
	return result, nil
}

func (r *roomRepo) AddRoomMember(ctx context.Context, member *biz.RoomMember, room *biz.Room) error {
	memberModel := &RoomMember{
		Id:       r.data.snowflake.Generate().Int64(),
		RoomId:   member.RoomId,
		UserId:   member.UserId,
		Role:     int(member.Role),
		JoinedAt: member.JoinedAt,
	}
	err := r.data.db.Transaction(func(tx *gorm.DB) error {
		var count int64 = 0
		err := tx.Model(&RoomMember{}).
			Where("room_id = ?", member.RoomId).
			WithContext(ctx).
			Count(&count).
			Error
		if err != nil {
			return err
		}
		if int(count) == room.MemberLimit {
			return biz.ErrMemberLimitReached
		}

		return tx.Model(memberModel).
			Create(&memberModel).
			WithContext(ctx).
			Error
	})
	return err
}

func (r *roomRepo) DeleteRoom(ctx context.Context, roomId int64) error {
	session, err := r.GetRoomSession(ctx, roomId)
	if err != nil {
		return err
	}

	return r.data.db.Transaction(func(tx *gorm.DB) error {
		if session != nil {
			conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(session.Endpoint))
			if err != nil {
				return err
			}
			defer conn.Close()
			_, err = gamingAPI.NewGamingClient(conn).DeleteGameInstance(ctx, &gamingAPI.DeleteGameInstanceRequest{
				RoomId: roomId,
				Force:  false,
			})
			if err != nil {
				return err
			}
			err = r.DeleteRoomSession(ctx, roomId, "")
			if err != nil {
				return err
			}
		}
		err = tx.Where("room_id = ?", roomId).Delete(&RoomMember{}).WithContext(ctx).Error
		if err != nil {
			return err
		}
		err = tx.
			Where("id = ?", roomId).
			Delete(&Room{}).
			WithContext(ctx).
			Error
		if err != nil {
			return err
		}
		err = cache.Del(ctx, r.data.rdb, roomCacheKey(roomId))
		if err != nil {
			r.log.Errorf("cache del room error: %v", err)
		}
		return nil
	})

}

func (r *roomRepo) UpdateRoom(ctx context.Context, room *biz.Room) error {
	err := r.data.db.Model(&Room{}).
		Where("id=?", room.Id).
		Updates(map[string]interface{}{
			"name":          room.Name,
			"private":       room.Private,
			"password_real": room.Password,
			"password_hash": room.PasswordHash,
		}).
		WithContext(ctx).
		Error
	if err != nil {
		return err
	}
	err = cache.Del(ctx, r.data.rdb, roomCacheKey(room.Id))
	if err != nil {
		r.log.Errorf("del room cache error: %v", err)
	}
	return nil
}

func (r *roomRepo) UpdateMember(ctx context.Context, member *biz.RoomMember) error {
	return r.data.db.Model(&RoomMember{}).
		Where("room_id = ? AND user_id = ?", member.RoomId, member.UserId).
		Updates(map[string]interface{}{
			"role": int(member.Role),
		}).
		WithContext(ctx).
		Error
}

func (r *roomRepo) DeleteMember(ctx context.Context, roomId, userId int64) error {
	return r.data.db.Debug().Where("room_id = ? AND user_id = ?", roomId, userId).
		Delete(&RoomMember{}).
		WithContext(ctx).
		Error
}

func (m *RoomMember) ToBizRoomMember() *biz.RoomMember {
	return &biz.RoomMember{
		UserId:   m.UserId,
		RoomId:   m.RoomId,
		Role:     v1.RoomRole(m.Role),
		Id:       m.Id,
		JoinedAt: m.JoinedAt,
	}
}

func (r *Room) ToBizRoom() *biz.Room {
	return &biz.Room{
		Id:           r.Id,
		Name:         r.Name,
		Host:         r.Host,
		Private:      r.Private,
		Password:     r.PasswordReal,
		PasswordHash: r.PasswordHash,
		MemberLimit:  r.MemberLimit,
		CreateTime:   r.CreatedAt,
		EmulatorType: r.EmulatorName,
	}
}

func (jr *JoinedRoom) ToBizJoinedRoom() *biz.JoinedRoom {
	return &biz.JoinedRoom{
		Room: biz.Room{
			Id:           jr.RoomId,
			Name:         jr.Name,
			Host:         jr.Host,
			Private:      jr.Private,
			Password:     jr.PasswordReal,
			PasswordHash: jr.PasswordHash,
			MemberLimit:  jr.MemberLimit,
			CreateTime:   jr.CreatedAt,
			EmulatorType: jr.EmulatorName,
		},
		UserId: jr.UserId,
		Role:   v1.RoomRole(jr.Role),
	}
}
