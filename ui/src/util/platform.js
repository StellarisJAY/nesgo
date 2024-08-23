const platform = {
    isMobile: function() {
        return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
    },
    isPortraitOrientation: function() {
        return this.isMobile &&
            (window.screen.orientation.type === "portrait-primary" ||
            window.screen.orientation.type==="portrait-secondary");
    },
}

export default platform;