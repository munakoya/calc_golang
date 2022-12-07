package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"os"
)

type Operand struct {
	// formからの受け取りはstringのみ
	// 小文字で始まるフィールド名は外部パッケージからのアクセス不可
	Title  string
	Op1    string
	Op2    string
	Op     string
	Result string
}

// テンプレートの設定 (サンプルコード)
// template構造体のポインタを返す
var tpl = template.Must(template.ParseFiles("index.html"))

// httpリクエストの処理
// w → サーバからクライアントへのレスポンス送信
// r → クライアントから送信されるリクエストメッセージ
func (i *Operand) mainHandler(w http.ResponseWriter, r *http.Request) {
	// httpヘッダの設定
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// フォームの入力値取得
	i.Op1 = r.FormValue("op1")
	i.Op2 = r.FormValue("op2")
	i.Op = r.FormValue("op") //演算子（ラジオボタンの値）

	// https://qiita.com/micropig3402/items/fae2b51e0a1632f9d796
	// math/bigパケ
	// big.int → cpuの影響を受けない整数型 → 32,64ビットを超えた大きい数字を取り扱える
	// string → Intに変換 (文字列 → 数値)
	convOp1 := &big.Int{}
	convOp2 := &big.Int{}

	// SetString(str, int)関数
	// 第1引数の文字列(数字)を第2引数で与えられた数値で解釈できればtrueを返す
	_, op1OK := convOp1.SetString(i.Op1, 10)
	_, op2OK := convOp2.SetString(i.Op2, 10)

	// 戻り値 → big.Int型に変換された数値とtrue or false

	if op1OK && op2OK {
		// big.Intパッケージ内にある四則演算できる関数を使用？
		resultInt := &big.Int{}
		switch i.Op {
		case "add":
			resultInt.Add(convOp1, convOp2)
		case "sub":
			resultInt.Sub(convOp1, convOp2)
		case "multi":
			resultInt.Mul(convOp1, convOp2)
		case "div":
			resultInt.Div(convOp1, convOp2)
		}
		// 計算結果の数値を文字列に変換 → htmlに文字列しか送れない？
		i.Result = resultInt.String()
	}

	// execute → テンプレートと構造体を組み合わせてhtmlクライアントへ送信
	outputHtml := tpl.Execute(w, i)

	// エラーがあればpanic()で終了
	if outputHtml != nil {
		panic(outputHtml)
	}
}

// https://www.morelife.work/entry/2019/09/14/data-uri-scheme-in-go-or-in-java
// html上に画像を表示 → Data URIスキームっていう方法で画像表示
// → 画像データをbase64のデータに変換
func imgOutput(imgPath string) {
	var err error
	file, err := os.Open(imgPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fi, _ := file.Stat()
	size := fi.Size()
	data := make([]byte, size)

	if _, err = file.Read(data); err != nil {
		log.Fatal(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	fmt.Println()
}
func main() {
	//ローカルサーバ立てる
	calc := &Operand{"電卓アプリ", "", "", "", ""}

	imgOutput("img/add.png")
	imgOutput("img/sub.png")
	imgOutput("img/mul.png")
	imgOutput("img/div.png")
	// ハンドラーcalcの設定
	http.HandleFunc("/calc", calc.mainHandler)

	// http.ListenAndServe関数 → httpサーバー起動
	// http.ListenAndServe(サーバアドレス, ルーティングハンドラ)

	// ルーティングハンドラが設定されているのならサーバ起動？
	result := http.ListenAndServe(":8080", nil)
	if result != nil {
		fmt.Println(result)
	}

}

// 苦労点 → htmlファイルに分離されていない → 分ける作業
// 拡張 → オペランドの数増やす → add関数とか使えなくなるかも？
