import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { AngularMaterialModule } from 'src/app/angular-material.module';
import { MatSnackBar } from '@angular/material/snack-bar';
import { GameComponent } from './game.component';
import { HttpClient } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { SidebarComponent } from './sidebar/sidebar.component';
import { NO_ERRORS_SCHEMA } from '@angular/compiler';
describe('GameComponent', () => {
    let component: GameComponent;
    let fixture: ComponentFixture<GameComponent>;
    let httpClient: HttpClient;
    let httpTestingController: HttpTestingController;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            declarations: [GameComponent, SidebarComponent],
            imports: [ HttpClientTestingModule, AngularMaterialModule, BrowserAnimationsModule ],
            providers: [ MatSnackBar, HttpClientTestingModule ],
            schemas: [ NO_ERRORS_SCHEMA ]
        }).compileComponents();

        // Inject the http service and test controller for each test
        httpClient = TestBed.get(HttpClient);
        httpTestingController = TestBed.get(HttpTestingController);

        fixture = TestBed.createComponent(GameComponent);
        component = fixture.componentInstance;
        component.ngOnInit();
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
