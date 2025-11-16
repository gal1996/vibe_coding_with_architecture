#!/bin/bash

# measure_tools.sh - ディレクトリ間の差分行数（LOC_churn）を計測するツール
#
# 使用方法:
#   ./measure_tools.sh <path1> <path2>
#
# 説明:
#   指定された2つのファイルまたはディレクトリ間の変更行数を集計します。
#   追加行数、削除行数、変更行数の合計を出力します。
#
# 出力形式:
#   Total LOC churn: XXX lines
#   - Added: XXX lines
#   - Deleted: XXX lines

# 引数チェック
if [ $# -ne 2 ]; then
    echo "エラー: 2つのパスを指定してください"
    echo "使用方法: $0 <path1> <path2>"
    echo ""
    echo "例:"
    echo "  $0 ./dir1 ./dir2          # ディレクトリ間の比較"
    echo "  $0 ./file1.txt ./file2.txt # ファイル間の比較"
    exit 1
fi

PATH1="$1"
PATH2="$2"

# パスの存在確認
if [ ! -e "$PATH1" ]; then
    echo "エラー: '$PATH1' が存在しません"
    exit 1
fi

if [ ! -e "$PATH2" ]; then
    echo "エラー: '$PATH2' が存在しません"
    exit 1
fi

# 一時ディレクトリの作成（ディレクトリ比較用）
TEMP_DIR=""
if [ -d "$PATH1" ] && [ -d "$PATH2" ]; then
    TEMP_DIR=$(mktemp -d)
    trap "rm -rf $TEMP_DIR" EXIT
fi

# diff実行と結果の解析
calculate_loc_churn() {
    local path1="$1"
    local path2="$2"

    # diffの実行（unified形式で出力）
    # -u: unified形式, -r: 再帰的（ディレクトリの場合）
    if [ -d "$path1" ] && [ -d "$path2" ]; then
        # ディレクトリ比較
        diff_output=$(diff -ur "$path1" "$path2" 2>/dev/null || true)
    else
        # ファイル比較
        diff_output=$(diff -u "$path1" "$path2" 2>/dev/null || true)
    fi

    # 追加行と削除行をカウント
    added_lines=$(echo "$diff_output" | grep -c "^+" | grep -v "^+++" || echo "0")
    deleted_lines=$(echo "$diff_output" | grep -c "^-" | grep -v "^---" || echo "0")

    # バイナリファイルの検出
    binary_files=$(echo "$diff_output" | grep -c "Binary files" || echo "0")

    echo "$added_lines $deleted_lines $binary_files"
}

# メイン処理
echo "========================================"
echo "LOC Churn 計測ツール"
echo "========================================"
echo ""
echo "比較対象:"
echo "  Path 1: $PATH1"
echo "  Path 2: $PATH2"
echo ""

# 計測実行
result=$(calculate_loc_churn "$PATH1" "$PATH2")
added=$(echo $result | cut -d' ' -f1)
deleted=$(echo $result | cut -d' ' -f2)
binary=$(echo $result | cut -d' ' -f3)

# 合計計算（変更行数 = 追加行数 + 削除行数）
total=$((added + deleted))

# 結果出力
echo "----------------------------------------"
echo "計測結果:"
echo "----------------------------------------"
echo "Total LOC churn: $total lines"
echo "  - Added:   $added lines"
echo "  - Deleted: $deleted lines"

if [ "$binary" -gt 0 ]; then
    echo ""
    echo "注意: $binary 個のバイナリファイルが検出されました（計測対象外）"
fi

echo ""
echo "========================================"

# CSV形式でも出力（オプション）
echo ""
echo "CSV形式（コピー用）:"
echo "$total,$added,$deleted"

# 終了コード
exit 0