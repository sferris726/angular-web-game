import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginButtonComponent } from './login-button/login-button.component';
import { provideHttpClient, withFetch } from '@angular/common/http';
import { RegisterButtonComponent } from './register-button/register-button.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginButtonComponent,
    RegisterButtonComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule
  ],
  providers: [provideHttpClient(withFetch())],
  bootstrap: [AppComponent]
})
export class AppModule { }
