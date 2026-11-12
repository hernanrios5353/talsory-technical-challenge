import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface QRResponse {
  q: number[][];
  r: number[][];
  statistics: {
    q: { max: number; min: number; average: number; sum: number; isDiagonal: boolean };
    r: { max: number; min: number; average: number; sum: number; isDiagonal: boolean };
    global: { max: number; min: number; average: number; sum: number; isDiagonal: boolean };
  };
}

@Injectable({ providedIn: 'root' })
export class QrService {
  // private apiUrl = 'https://talsory-api-go.onrender.com/api/v1';
  private apiUrl = '/api/v1'; // Cambiado a localhost para pruebas locales

  constructor(private http: HttpClient) {}

  login(username: string, password: string): Observable<{ token: string }> {
    return this.http.post<{ token: string }>(`${this.apiUrl}/auth/login`, { username, password });
  }

  processMatrix(matrix: number[][], token: string): Observable<QRResponse> {
    const headers = { Authorization: `Bearer ${token}` };
    return this.http.post<QRResponse>(`${this.apiUrl}/qr`, { matrix }, { headers });
  }
}