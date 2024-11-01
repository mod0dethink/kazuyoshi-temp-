package question

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func MakeQuestionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ユーザー認証の確認
		cookie, err := r.Cookie("username")
		if err != nil {
			http.Error(w, "ログインが必要です", http.StatusUnauthorized)
			return
		}

		// リクエストボディの解析
		var question Question
		if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
			http.Error(w, "無効なリクエストデータです", http.StatusBadRequest)
			return
		}

		// バリデーション
		if len(question.Choices) != 4 {
			http.Error(w, "選択肢は4つ必要です", http.StatusBadRequest)
			return
		}

		// データベースに問題を保存
		_, err = db.Exec(
			"INSERT INTO questions (creator_username, question_text, correct_answer, choice1, choice2, choice3, choice4, explanation) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			cookie.Value,
			question.QuestionText,
			question.CorrectAnswer,
			question.Choices[0],
			question.Choices[1],
			question.Choices[2],
			question.Choices[3],
			question.Explanation,
		)

		if err != nil {
			http.Error(w, "問題の保存に失敗しました", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "問題が正常に作成されました",
		})
	}
}
