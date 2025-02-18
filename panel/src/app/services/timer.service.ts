import { Injectable } from '@angular/core';
import { interval, Subscription } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TimerService {
  private subscription: Subscription = new Subscription;

  startTimer(intervalInMinutes: number, callback: () => void) {
    this.subscription = interval(intervalInMinutes * 60 * 1000)
      .subscribe(() => {
        if (callback) {
          callback()
        }
      })
  }

  stopTimer() {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }
}
