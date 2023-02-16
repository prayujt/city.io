import { Component, OnInit } from '@angular/core';
import { CookieService } from 'ngx-cookie-service';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environments/environment';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css'],
})
export class AppComponent implements OnInit{
    constructor(
        private http: HttpClient,
        private cookieService: CookieService,
        private _snackBar: MatSnackBar,
        private router: Router,
    ) {}
    private ID: string = '';
    public loggedIn: boolean = false;

    public ngOnInit(): void {
        setInterval(() => {
            this.ID = this.cookieService.get('cookie');
            if (this.ID != "") {
                this.http.get<any>(
                    `http://${environment.API_HOST}:${environment.API_PORT}/sessions/${this.ID}`
                )
                .subscribe((response) => {
                    if (response) { 
                        this.loggedIn = true 
                    }
                })
            }
        }, 250);
    }

    public logOut(): void {
        this.http.post<any>(
            `http://${environment.API_HOST}:${environment.API_PORT}/sessions/logout`,
            {
                sessionId: this.ID
            })
            .subscribe((response) => {
                if (response) {
                    this._snackBar.open(
                        'Log out successful!',
                        'Close',
                        {
                            duration: 2000,
                        }
                    );
                    this.loggedIn = false;
                } else {
                    this._snackBar.open(
                        'Could not log out!',
                        'Close',
                        {
                            duration: 2000,
                        }
                    );
                }
            })
    }

    title = 'client';

}
