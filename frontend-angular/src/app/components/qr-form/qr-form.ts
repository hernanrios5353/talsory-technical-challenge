import { Component, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { QrService, QRResponse } from '../../services/qr';
import { finalize } from 'rxjs/operators';

@Component({
  selector: 'app-qr-form',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './qr-form.html',
  styleUrl: './qr-form.scss'
})
export class QrFormComponent {
  username = 'admin';
  password = 'password123';
  token = '';
  matrixJson = '[[1,2],[3,4],[5,6]]';
  result: QRResponse | null = null;
  errorMessage = '';
  isLoadingLogin = false;
  isLoadingProcess = false;

  // Inyectamos ChangeDetectorRef
  constructor(private qrService: QrService, private cdr: ChangeDetectorRef) {}

  login() {
    this.errorMessage = '';
    this.isLoadingLogin = true;
    
    this.qrService.login(this.username, this.password).subscribe({
      next: (res) => {
        this.token = res.token;
        this.isLoadingLogin = false;
        this.cdr.detectChanges(); // <-- Forzamos la actualización de la vista
      },
      error: (err) => {
        console.error(err);
        this.errorMessage = 'Error en login. Verifica la conexión.';
        this.isLoadingLogin = false;
        this.cdr.detectChanges(); // <-- Forzamos la actualización
      }
    });
  }

  process() {
    this.errorMessage = '';
    this.isLoadingProcess = true;
    
    try {
      const matrix = JSON.parse(this.matrixJson);
      this.qrService.processMatrix(matrix, this.token).subscribe({
        next: (res) => {
          this.result = res;
          this.isLoadingProcess = false;
          this.cdr.detectChanges(); // <-- Forzamos la actualización
        },
        error: (err) => {
          console.error(err);
          this.errorMessage = 'Error en cálculo.';
          this.isLoadingProcess = false;
          this.cdr.detectChanges(); // <-- Forzamos la actualización
        }
      });
    } catch (e) {
      this.errorMessage = 'JSON inválido.';
      this.isLoadingProcess = false;
      this.cdr.detectChanges(); // <-- Forzamos la actualización
    }
  }

  formatNumber(num: number): string {
    return Number(num.toFixed(4)).toString();
  }
}