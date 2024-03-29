import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';

import { ConstructableService } from './constructable.service';

describe('ConstructableService', () => {
    let service: ConstructableService;

    beforeEach(() => {
        TestBed.configureTestingModule({
            imports: [HttpClientTestingModule],
            providers: [HttpClientTestingModule],
        });
        service = TestBed.inject(ConstructableService);
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });
});
