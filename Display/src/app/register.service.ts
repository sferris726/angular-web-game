import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class RegisterService {
  constructor(private _http: HttpClient) {}

  public register(
    username: string,
    email: string,
    pw: string,
    settings_box: number,
    settings_guess: number
  ) {
    const url = 'http://localhost:3000/register';
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    this._http
      .post(
        url,
        {
          username: username,
          email: email,
          password: pw,
          settings_box: settings_box,
          settings_guess: settings_guess,
        },
        { headers: headers }
      )
      .subscribe((res) => {
        console.log('RECEIVED ', res);
      });
  }
}
