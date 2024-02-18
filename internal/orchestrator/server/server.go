package server

import (
	"calculator/internal/model"
	expressionParser "calculator/internal/orchestrator/expressionParser"
	"calculator/pkg/rabbitmq/orchestrator"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Orchestrator struct {
	mux *mux.Router
}

type MyRequest struct {
	Expression string `json:"expression"`
	Id         string `json:"-"`
}

var cache LRUCache

func New() Orchestrator {
	serveMux := mux.NewRouter()

	// Хэндлеры для запросов с сайта
	serveMux.HandleFunc("/", mainPage)

	serveMux.Handle("/expression", ValidateExpressionMiddleware(http.HandlerFunc(ExpressionHandler))).Methods("POST")

	serveMux.Handle("/agentstate", nil)

	serveMux.HandleFunc("/expression", GetExpressionById).Methods("GET")

	http.Handle("/", serveMux)

	// Хэндлеры для обработки запросов от агента

	// Хэндлеры для API
	serveMux.HandleFunc("/api/getimpodencekey", GetImpodenceKeyHandler).Methods("POST")
	return Orchestrator{serveMux}
}

func (orchestrator Orchestrator) Run() {
	// Создаём кэш
	cache = *NewLRUCache(30)

	// Добавляем middleware для перехвата паник
	muxWithPanicHandler := PanicMiddleware(orchestrator.mux)

	server := http.Server{
		Addr:    ":8080",
		Handler: muxWithPanicHandler,
	}

	server.ListenAndServe()
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, use curl :)")
}

func GetExpressionById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if el, ok := cache.Get(id); ok {
		fmt.Fprintf(w, "id: %v\n", el.IdExpression)
		fmt.Fprintf(w, "expression: %v\n", el.InfinixExpression)
		fmt.Fprintf(w, "answer: %v\n", el.Result)
		fmt.Fprintf(w, "status: %v\n", el.Status)
		fmt.Fprintf(w, "createdAt: %v\n", el.CreatedAt.String())
		fmt.Fprintf(w, "completedAt: %v\n", el.SolvedAt.String())
		return
	}

	fmt.Fprintln(w, "Expression doesn't exist")
}

func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	expression := r.Context().Value("expression").(model.Expression)

	// Проверяем выражение на наличие результата в кэше (если нет в кэше - кэш пойдёт спрашивать у базы данных, в ином случае отправляем агенту на вычисление)
	if el, ok := cache.Get(expression.IdExpression); ok {
		fmt.Fprintln(w, el.Result)
		return
	}

	expression, err := orchestrator.SolveExpression(&expression)
	if err != nil {
		log.Println(err.Error())
	}
	cache.Set(expression.IdExpression, expression)
	fmt.Fprintln(w, expression.Result)
}

// вообще по идее оно должно создаваться на фронтэнде, но т.к пока нет фронта - создаём на бэкенде (по запросу с фронта)
func GetImpodenceKeyHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим JSON
	var request MyRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Error while parsing JSON", http.StatusBadRequest)
		return
	}

	request.Expression = strings.ReplaceAll(request.Expression, " ", "")

	// и получаем ключ
	key := expressionParser.CreateImpodenceKey(request.Expression)

	w.Write([]byte(key))
}

// Проверяем на валидность наше выражение и передаём уже паршенное (парсенное?) в Handler
func ValidateExpressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Парсим JSON
		var request MyRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Error while parsing JSON", http.StatusBadRequest)
			return
		}

		// Парсим полученное выражение
		postfixExpression, err := expressionParser.ParseExpression(request.Expression)
		if err != nil {
			http.Error(w, "Error while parsing expression", http.StatusBadRequest)
			return
		}

		expression := model.Create(request.Expression, postfixExpression, expressionParser.CreateImpodenceKey(request.Expression))

		// Передаём в реквест контекст с выражением для дальнейшей работы с ним в хэндлере
		rWithContext := r.WithContext(context.WithValue(r.Context(), "expression", expression))
		// Пишем в хэдер статус код, что всё хорошо
		w.WriteHeader(http.StatusAccepted)
		next.ServeHTTP(w, rWithContext)
	})
}

// Мидлварь для отлавливания паник
func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
