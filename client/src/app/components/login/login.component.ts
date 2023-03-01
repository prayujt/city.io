import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';

import { environment } from '../../../environments/environment';
import { CookieService } from 'ngx-cookie-service';

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css'],
})
export class LoginComponent {
    public cookieValue: string = '';
    constructor(
        private http: HttpClient,
        private _snackBar: MatSnackBar,
        private router: Router,
        private cookieService: CookieService
    ) {}
    public ngOnInit(): void {
        let sessionId: string = this.cookieService.get('sessionId');
        if (sessionId != '') {
            this.http
                .get<any>(
                    `http://${environment.API_HOST}:${environment.API_PORT}/sessions/${sessionId}`
                )
                .subscribe((response) => {
                    if (response.status) {
                        this.router.navigate(['game']);
                    }
                });
        }
    }
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
                    `http://${environment.API_HOST}:${environment.API_PORT}/login/createSession`,
                    {
                        username: username,
                        password: password,
                    }
                )
                .subscribe((response) => {
                    console.log(response)
                    if (response.sessionId != '') {
                        this.router.navigate(['game']);
                        this.cookieService.set('sessionId', response.sessionId);
                        this.cookieValue = this.cookieService.get('sessionId');
                    } else {
                        this._snackBar.open(
                            'Invalid username or password!',
                            'Close',
                            {
                                duration: 2000,
                            }
                        );
                    }
                });
        }
    }
}
