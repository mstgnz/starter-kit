import { Component, Input, OnInit } from '@angular/core';
import * as XLSX from 'xlsx';
import { AlertifyService } from '../../services/alertify.service';

@Component({
  selector: 'app-shared-excel-export',
  standalone: true,
  templateUrl: './excel-export.component.html',
  styles: ``
})
export class ExcelExportComponent implements OnInit {

  @Input() document: any[] = []

  constructor(
    private alertifyService: AlertifyService
  ) { }

  ngOnInit(): void { }

  fileExcelExport() {

    if (this.document.length) {
      // export
      const ws: XLSX.WorkSheet = XLSX.utils.json_to_sheet(this.document);

      /* generate workbook and add the worksheet */
      const wb: XLSX.WorkBook = XLSX.utils.book_new();
      XLSX.utils.book_append_sheet(wb, ws, 'Export');

      /* save to file */
      XLSX.writeFile(wb, 'tursys-export.xlsx');
    } else {
      this.alertifyService.error('HATA!, Boş data export edilemez')
    }

  }

}
