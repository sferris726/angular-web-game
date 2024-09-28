import { Injectable } from "@angular/core";
import { HttpClient, HttpHeaders } from "@angular/common/http";

@Injectable({
  providedIn: "root",
})
export class LoginService {
  constructor(private _http: HttpClient) {}

  public sendCredentials(username: string, pw: string) {
    const url = "http://localhost:3000/login";
    const headers = new HttpHeaders({ "Content-Type": "application/json" });
    return this._http.post(url, { username, pw }, { headers: headers });
  }
}
