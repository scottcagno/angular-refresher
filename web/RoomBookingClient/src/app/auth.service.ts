import {EventEmitter, Injectable} from '@angular/core';
import {DataService} from "./data.service";

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  isAuthenticated = false;
  authResultEvent = new EventEmitter<boolean>();

  constructor(private dataService : DataService) { }

  authenticate(username : string, password : string) {
    this.dataService.validateUser(username, password).subscribe(
      next => {
        this.isAuthenticated = true;
        this.authResultEvent.emit(true);
      },
      error => {
        this.isAuthenticated = false;
        this.authResultEvent.emit(false);
      }
    );
  }

  logout() {
    this.isAuthenticated = false;
  }
}
