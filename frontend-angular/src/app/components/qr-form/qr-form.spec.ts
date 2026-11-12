import { ComponentFixture, TestBed } from '@angular/core/testing';
import { QrForm } from './qr-form';

describe('QrForm', () => {
  let component: QrForm;
  let fixture: ComponentFixture<QrForm>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [QrForm],
    }).compileComponents();

    fixture = TestBed.createComponent(QrForm);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
