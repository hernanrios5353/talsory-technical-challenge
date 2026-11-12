import { Component } from '@angular/core';
import { QrFormComponent } from './components/qr-form/qr-form';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [QrFormComponent],
  template: `<app-qr-form></app-qr-form>`,
})
export class App {}