import { Component, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';

@Component({
    selector: 'app-register',
    templateUrl: './register.component.html',
    styleUrls: ['./register.component.css'],
})
@Injectable()
export class RegisterComponent {
    constructor(
        private http: HttpClient,
        private _snackBar: MatSnackBar,
        private router: Router
    ) {}
    ngOnInit() {}
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
            console.log('Valid Password');
            this.http
                .post<any>('http://prayujt.com:8000/login/createAccount', {
                    username: username,
                    password: password,
                })
                .subscribe((response) => {
                    console.log(response);
                    if (!response) {
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
