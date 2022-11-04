import {Component, EventEmitter, Input, OnDestroy, OnInit, Output} from '@angular/core';
import {User} from "../../../model/User";
import {DataService} from "../../../data.service";
import {Router} from "@angular/router";
import {FormResetService} from "../../../form-reset.service";
import {Subscription} from "rxjs";

@Component({
  selector: 'app-user-edit',
  templateUrl: './user-edit.component.html',
  styleUrls: ['./user-edit.component.css']
})
export class UserEditComponent implements OnInit, OnDestroy {

  @Input()
  user !:User;

  @Output()
  dataChangedEvent = new EventEmitter();

  formUser !:User; // just for the form
  password !:string;
  password2 !:string;
  message !:string;
  nameIsValid = false;
  passwordsAreValid = false;
  passwordsMatch = false;
  userFormReset !: Subscription;

  constructor(private dataService :DataService,
              private router :Router,
              private formResetService: FormResetService) {
  }

  ngOnInit(): void {
    this.initializeForm();
    this.userFormReset = this.formResetService.resetUserFormEvent.subscribe(
      user => {
        this.user = user;
        this.initializeForm();
      }
    );
  }

  ngOnDestroy() {
    this.userFormReset.unsubscribe();
  }

  initializeForm() {
    this.formUser = Object.assign({}, this.user);
    this.checkIfNameIsValid();
    this.checkIfPasswordsAreValid();
  }

  onCancel() {
    this.router.navigate(['admin','users'], {queryParams:{id: this.user.id, action:'view'}});
  }

  onSubmit() {
    // check add new or edit existing
    if (this.formUser.id == null) {
      // add new
      this.dataService.addNewUser(this.formUser, this.password).subscribe((user)=>{
        this.dataChangedEvent.emit();
        this.router.navigate(['admin','users'], {queryParams:{id: user.id, action: 'view'}});
      },
        error => {
          this.message = 'Something went wrong and the data was not saved. You may want to try again.';
        }
    );
    } else {
      // edit existing
      this.dataService.updateUser(this.formUser).subscribe((user)=>{
        this.dataChangedEvent.emit();
        this.router.navigate(['admin','users'], {queryParams:{id: user.id, action: 'view'}});
      },
        error => {
        this.message = 'Something went wrong and the data was not saved. You may want to try again.';
        });
    }
  }

  checkIfNameIsValid() {
    if (this.formUser.name) {
      this.nameIsValid = this.formUser.name.trim().length > 0;
    } else {
      this.nameIsValid = false;
    }
  }

  checkIfPasswordsAreValid() {
    if (this.formUser.id != null) {
      this.passwordsAreValid = true;
      this.passwordsMatch = true;
    } else {
      this.passwordsMatch = this.password === this.password2;
      if (this.password) {
        this.passwordsAreValid = this.password.trim().length > 0;
      } else {
        this.passwordsMatch = false;
      }
    }
  }

}
