import { Injectable } from "@angular/core";
import * as XLSX from 'xlsx';
import { AlertifyService } from "./alertify.service";
import { DateTimeColumn } from "../interfaces/select.interface";

@Injectable()
export class ExcelService {

    constructor(
        private alertifyService: AlertifyService
    ) { }

    import(file: any, dateTimeColumn?: DateTimeColumn) {
        return new Promise<any[]>(resolve => {
            const target: DataTransfer = <DataTransfer>(file.target)
            if (target.files.length !== 1) throw new Error('Cannot use multiple files')
            const reader: FileReader = new FileReader()
            reader.readAsBinaryString(target.files[0])
            reader.onload = (e: any) => {
                const bstr: string = e.target.result
                const wb: XLSX.WorkBook = XLSX.read(bstr, { type: 'binary' })
                const wsname: string = wb.SheetNames[0]
                const ws: XLSX.WorkSheet = wb.Sheets[wsname]
                resolve(this.parseExcelSheet(ws, dateTimeColumn))
            }
        })
    }

    export(document: any, name: string = "export") {
        return new Promise<boolean>(resolve => {
            if (document.length) {
                // export
                const ws: XLSX.WorkSheet = XLSX.utils.json_to_sheet(document);
                /* generate workbook and add the worksheet */
                const wb: XLSX.WorkBook = XLSX.utils.book_new();
                XLSX.utils.book_append_sheet(wb, ws, name);
                /* save to file */
                XLSX.writeFile(wb, `tursys-${name}.xlsx`);
                resolve(true)
            } else {
                this.alertifyService.error('HATA!, Boş data export edilemez')
                resolve(true)
            }
        })
    }

    parseExcelSheet(ws: XLSX.WorkSheet, dateTimeColumn?: DateTimeColumn): any[] {
        const data: any[] = [];
        XLSX.utils.sheet_to_json(ws).forEach((row: any) => {
            for (const col in row) {
                if (col.endsWith("id")) continue
                if (row.hasOwnProperty(col) && typeof row[col] === 'number') {
                    const cellValue = row[col]
                    const parseDate = XLSX.SSF.parse_date_code(cellValue)
                    if (parseDate && parseDate.y >= 2000) {
                        const formattedDate = XLSX.SSF.format('yyyy-mm-dd hh:mm', cellValue)
                        row[col] = formattedDate
                    }
                }
            }
            if (dateTimeColumn) {
                if (dateTimeColumn.time) {
                    dateTimeColumn.time.forEach((tc: string) => {
                        if (row[tc]) {
                            row[tc] = XLSX.SSF.format('hh:mm', row[tc])
                        }
                        if (row[tc] === 0) {
                            row[tc] = "00:00"
                        }
                    })
                }
                if (dateTimeColumn.date) {
                    dateTimeColumn.date.forEach((dc: string) => {
                        if (row[dc]) {
                            row[dc] = XLSX.SSF.format('yyyy-mm-dd hh:mm', row[dc])
                        }
                    })
                }
            }
            data.push(row)
        })
        return data
    }

}