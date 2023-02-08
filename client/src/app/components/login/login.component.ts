import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';

import { environment } from '../../../environments/environment';

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css'],
})
export class LoginComponent {
    constructor(
        private http: HttpClient, 
        private _snackBar: MatSnackBar, 
        private router: Router
    ) {}
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
            this.http
                .post<any>(
                    `http://${environment.API_HOST}:${environment.API_PORT}/login/verifyAccount`,
                    {
                        username: username,
                        password: password,
                    }
                )
                .subscribe((response) => {
                    if (response['status']) {
                        this.router.navigate(['game']);
                    }
                    else {
                        this._snackBar.open('Invalid username or password!', 'Close', {
                            duration: 2000,
                        });
                    }
                });
        }
    }
}
