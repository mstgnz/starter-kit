/// <reference lib="dom" />

import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import * as XLSX from 'xlsx';

@Component({
  selector: 'app-shared-excel-import',
  standalone: true,
  templateUrl: './excel-import.component.html',
  styles: ``
})
export class ExcelImportComponent implements OnInit {

  @Output() document = new EventEmitter<any>()

  constructor() { }

  ngOnInit(): void { }

  openUpload() {
    document.getElementById('excel-import')?.click()
  }

  fileExcelUpload(file: any) {
    /* wire up file reader */
    const target: DataTransfer = <DataTransfer>(file.target)
    if (target.files.length !== 1) throw new Error('Cannot use multiple files')
    const reader: FileReader = new FileReader()
    reader.readAsArrayBuffer(target.files[0])
    reader.onload = (e: any) => {
      /* read workbook */
      const bstr: string = e.target.result
      const wb: XLSX.WorkBook = XLSX.read(bstr, { type: 'binary' })
      /* grab first sheet */
      const wsname: string = wb.SheetNames[0]
      const ws: XLSX.WorkSheet = wb.Sheets[wsname]
      /* save data */
      //const data = <any>(XLSX.utils.sheet_to_json(ws, { header: 1 }))
      //data.splice(0, 1)
      const data = <any>(XLSX.utils.sheet_to_json(ws))
      this.document.emit(data)
    }
    file.target.value = ''
  }

}
