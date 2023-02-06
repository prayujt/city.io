import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css'],
})
export class LoginComponent {
    constructor(private http: HttpClient, private _snackBar: MatSnackBar) {}
    public showPassword: boolean = false;
    public toggleVisibility(): void {
        this.showPassword = !this.showPassword;
    }
    public getData(username: string, password: string): void {
        if (username.length == 0) {
            this._snackBar.open('Input a valid username!', 'Close', {
                duration: 2000,
            });
        } else if (password.length == 0) {
            this._snackBar.open('Input a valid password!', 'Close', {
                duration: 2000,
            });
        } else {
            console.log('Valid Password');
            this.http
                .post<any>('http://prayujt.com:8000/login/verifyAccount', {
                    username: username,
                    password: password,
                })
                .subscribe((response) => {
                    console.log(response);
                });
        }
    }
}
