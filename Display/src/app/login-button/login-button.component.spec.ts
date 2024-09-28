import { TestBed } from '@angular/core/testing';
import { LoginButtonComponent } from './login-button.component';

describe('LoginButtonComponent', () => {
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LoginButtonComponent],
    }).compileComponents();
  });

  it('should create the app', () => {
    const fixture = TestBed.createComponent(LoginButtonComponent);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });

  it(`should have the 'my-app' title`, () => {
    const fixture = TestBed.createComponent(LoginButtonComponent);
    const app = fixture.componentInstance;
  });

  it('should render title', () => {
    const fixture = TestBed.createComponent(LoginButtonComponent);
    fixture.detectChanges();
  });
});
