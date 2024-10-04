import { Component } from '@angular/core';
import { RegisterService } from '../register.service';

@Component({
  selector: 'app-register-button',
  templateUrl: './register-button.component.html',
  styleUrl: './register-button.component.scss',
})
export class RegisterButtonComponent {
  constructor(private _register_service: RegisterService) {}

  onSubmit() {
    this._register_service.register(
      'Scott',
      'Scott@gmail.com',
      'helloworld',
      10,
      5
    );
  }
}
