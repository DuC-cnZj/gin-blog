/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/youngduc/go-blog/cmd"
)

//func timeMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
//
//		// next handler
//		next.ServeHTTP(wr, r)
//
//		log.Println(r.Response)
//	})
//}

func main() {
	//http.Handle("/post", timeMiddleware(H()))
	//http.ListenAndServe(":8989", nil)
	cmd.Execute()
}

//func H()  http.Handler{
//	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
//		wr.WriteHeader(200)
//		wr.Write([]byte("ok"))
//	})
//}
