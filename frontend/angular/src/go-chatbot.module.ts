import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { GoChatbotComponent } from './go-chatbot.component';

@NgModule({
  declarations: [
    GoChatbotComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    HttpClientModule
  ],
  exports: [
    GoChatbotComponent
  ]
})
export class GoChatbotModule { }
