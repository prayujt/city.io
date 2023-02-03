import { Component } from '@angular/core';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  public showPassword: boolean = false;
  public toggleVisibility(): void {
    this.showPassword = !this.showPassword;
  }
}
