import { Component } from '@angular/core';
import { LoginService } from '../login.service';

@Component({
  selector: 'app-login-button',
  templateUrl: './login-button.component.html',
  styleUrl: './login-button.component.scss'
})
export class LoginButtonComponent {

  constructor(private _login_service: LoginService) {
  }

  public onSubmit() {
    console.log("WE IN HERE BUTTON");
    this._login_service.sendCredentials("Scott", "1234");
  }
}
