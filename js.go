package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/gopherjs/gopherjs/js"
	"strings"
)

//ループの一周目かどうか
var isFirstJS = true

//ループの中の処理
func jsEvent() {
	//ループの最初
	if isFirstJS {
		fitScreen()
		isFirstJS = false
	}

	//リサイズ時
	js.Global.Call("addEventListener", "resize", func() {
		fitScreen()
	})
}

func fitScreen() {
	//スケールを計算
	scale := calcScale()
	//画面サイズを調整
	ebiten.SetScreenScale(scale)
}

//スケールを計算
func calcScale() float64 {
	//ウインドウの大きさを取得
	innerWidth := js.Global.Get("window").Get("innerWidth").Float()
	innerHeight := js.Global.Get("window").Get("innerHeight").Float()

	//ウインドウサイズ/キャンバスのサイズ
	scaleWidth := innerWidth / screenWidth	
	scaleHeight := innerHeight / screenHeight

	scale := 1.0

	//倍率の小さいほうを全体のスケールにする
	if scaleWidth < scaleHeight {
		scale = scaleWidth
	}else {
		scale = scaleHeight
	}

	return scale		
}

//デバイス(userAgent)を取得
func getDevice() bool {
	var isMobile bool
	ua := js.Global.Get("navigator").Get("userAgent").String()

	if strings.Index(ua, "iPhone") != -1 || strings.Index(ua, "iPod") != -1 || strings.Index(ua, "Android") != -1 || strings.Index(ua, "Mobile") != -1 {
		//スマホ
		isMobile = true
	}else {
		//スマホ以外
		isMobile = false
	}

	return isMobile
}



