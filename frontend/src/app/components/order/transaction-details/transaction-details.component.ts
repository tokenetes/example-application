import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { TransactionDetails } from '../../../models/transaction-details.model';

@Component({
  selector: 'app-transaction-details',
  templateUrl: './transaction-details.component.html',
  styleUrls: ['./transaction-details.component.css'],
  standalone: false
})
export class TransactionDetailsComponent implements OnInit {
  transactionDetails: TransactionDetails | null = null;

  constructor(private router: Router) { }

  ngOnInit(): void {
    const storedDetails = sessionStorage.getItem('transactionDetails');
    if (storedDetails) {
      this.transactionDetails = JSON.parse(storedDetails);
      sessionStorage.removeItem('transactionDetails');
    } else {
      console.error('No transaction details available');
      this.router.navigate(['/']);
    }
  }
}