import { Component, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})

@Injectable()
export class RegisterComponent {
  constructor(private http: HttpClient) {}
  ngOnInit() {}
  public getData(username: string, password: string, password2: string): void {
    if (password == password2){
      this.http.post<any>('http://prayujt.com:8000/login/createAccount', {
        username: username,
        password: password
      }).subscribe((response) => {
        console.log(response);
      })
    }
  }
}
