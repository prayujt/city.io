import { Component } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';

import { environment } from '../../../environments/environment';
import { CookieService } from 'ngx-cookie-service';

@Component({
    selector: 'app-register',
    templateUrl: './register.component.html',
    styleUrls: ['./register.component.css'],
})
export class RegisterComponent {
    constructor(
        private http: HttpClient,
        private _snackBar: MatSnackBar,
        private router: Router,
        private cookieService: CookieService
    ) {}

    public ngOnInit(): void {
        let jwtToken: string = this.cookieService.get('jwtToken');
        if (jwtToken != '') {
            let headers = new HttpHeaders();
            headers = headers.append('Token', jwtToken);
            this.http
                .get<any>(`${environment.API_HOST}/sessions/validate`, {
                    headers,
                })
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
    public getData(
        username: string,
        password: string,
        password2: string
    ): void {
        if (username.length == 0) {
            this._snackBar.open('Input a valid username!', 'Close', {
                duration: 2000,
            });
        } else if (password.length == 0) {
            this._snackBar.open('Input a valid password!', 'Close', {
                duration: 2000,
            });
        } else if (password != password2) {
            this._snackBar.open('Passwords must match!', 'Close', {
                duration: 2000,
            });
        } else {
            this.http
                .post<any>(`${environment.API_HOST}/login/createAccount`, {
                    username: username,
                    password: password,
                })
                .subscribe((response) => {
                    if (!response.status) {
                        this._snackBar.open(
                            'This username is not available!',
                            'Close',
                            {
                                duration: 2000,
                            }
                        );
                    } else {
                        this._snackBar.open(
                            'Account successfully created!',
                            'Close',
                            {
                                duration: 2000,
                            }
                        );
                        this.router.navigate(['login']);
                    }
                });
        }
    }
}
