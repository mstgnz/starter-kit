import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment.development';

@Injectable()
export class HashService {
    private secret: string = environment.appSecret

    public async setHeaderHash(url: string): Promise<{ Timestamp: string; Hash: string }> {
        const timestamp = Math.floor(Date.now() / 1000).toString();
        const rawData = `Saha.${timestamp}:${url}:${this.secret}.Kolay`;
        const hash = await this.sha256(rawData);

        return {
            Timestamp: timestamp,
            Hash: hash,
        };
    }

    private async sha256(data: string): Promise<string> {
        const encoder = new TextEncoder();
        const dataBuffer = encoder.encode(data);
        const hashBuffer = await crypto.subtle.digest('SHA-256', dataBuffer);
        return Array.from(new Uint8Array(hashBuffer))
            .map((byte) => byte.toString(16).padStart(2, '0'))
            .join('')
            .toLowerCase();
    }
}
