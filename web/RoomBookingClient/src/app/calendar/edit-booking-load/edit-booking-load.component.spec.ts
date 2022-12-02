import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EditBookingLoadComponent } from './edit-booking-load.component';

describe('EditBookingLoadComponent', () => {
  let component: EditBookingLoadComponent;
  let fixture: ComponentFixture<EditBookingLoadComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ EditBookingLoadComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(EditBookingLoadComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
