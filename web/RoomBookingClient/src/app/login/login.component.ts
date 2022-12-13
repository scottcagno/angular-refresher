import {Component, OnDestroy, OnInit} from '@angular/core';
import {AuthService} from "../auth.service";
import {ActivatedRoute, Router} from "@angular/router";
import {Subscription} from "rxjs";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit, OnDestroy {

  message = '';

  username !: string;
  password !: string;
  subscription !: Subscription;

  constructor(private authService : AuthService,
              private route : Router,
              private activatedRoute : ActivatedRoute) { }

  ngOnInit(): void {
    this.subscription = this.authService.authResultEvent.subscribe(result => {
      if (result) {
        const url = this.activatedRoute.snapshot.queryParams['requested'];
        this.route.navigateByUrl(url);
      } else {
        this.message = 'Your username or password was not recognized, please try again.';
      }
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  onSubmit() {


    this.authService.authenticate(this.username, this.password)
  }

  onCancel() {
  }

  onLogout() {
    this.authService.logout();
    this.route.navigate(['']);
  }

}
