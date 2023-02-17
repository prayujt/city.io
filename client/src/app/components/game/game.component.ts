import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

import { environment } from '../../../environments/environment';
import { CookieService } from 'ngx-cookie-service';

@Component({
    selector: 'app-game',
    templateUrl: './game.component.html',
    styleUrls: ['./game.component.css'],
})
export class GameComponent {
    constructor(
        private http: HttpClient,
        private router: Router,
        private cookieService: CookieService
    ) {}
    public ngOnInit(): void {
        this.getID();
        let ID: string = this.cookieService.get('cookie');
        if (ID != '') {
            this.http
                .get<any>(
                    `http://${environment.API_HOST}:${environment.API_PORT}/sessions/${ID}`
                )
                .subscribe((response) => {
                    if (!response) {
                        this.router.navigate(['login']);
                    }
                });
        }
    }
    public getID(): string {
        console.log(this.cookieService.get('cookie'));
        return this.cookieService.get('cookie');
    }
}
