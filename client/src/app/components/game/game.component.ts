import { Component } from '@angular/core';
import { CookieService } from 'ngx-cookie-service';

@Component({
    selector: 'app-game',
    templateUrl: './game.component.html',
    styleUrls: ['./game.component.css'],
})
export class GameComponent {
    constructor(
        private cookieService: CookieService
    ) {}
    public ngOnInit(): void {
        this.getID();
    }
    public getID(): string {
        console.log(this.cookieService.get('cookie'));
        return this.cookieService.get('cookie');
    }
}
