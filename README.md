# Project Title

    AICULUS DASHBOARD

# Executive Summary

    AiCULUS dashboard is a web - based project which helps users to view graphical representation of index data for providing trade actions so that the users can easily analyse the data. 

# Contributors

    Name - PreferName- StudentID - EmailAddress

    <Binbin Wang - Leo - s3625685 - 418624658@qq.com>
    <Zhuojun Li - Elvin - s3514856 - 1169219351@qq.com>
    <Sethu Mohan Das - Sethu - s3630734 - sethumohandas@gmail.com>
    <Tun-Ta Hsieh - Darren - s3584414 - haha401136@gmail.com>
    <Fen Gan - Gary - s3703529 - ganfen9173@gmail.com>
    <Shuijia Zhuo - Roger - s3384039 - shuijia179@gmail.com>

# Description

### .idea & .vscode

    These are the launch&setting files of IDE.

### css

    The folder which contains css file.

### data

    The folder which contains data files. 

### img

    The folder which contains logos. 

### log

    The folder which contains log files.

### static

    Chart.js file.

### views

    Html files location.

### DAO.go

    The back-end file.

### main.go

    Main function to run this dashboard.

### routes.go

    The controller between front-end and back-end.

# Implementation  

## 1. GOLANG

    ```
    Please CHECK your GoLang system environment. 
    There should be a "GOPATH" with "C:\Go\bin" (your GO location), and a "GOROOT" with "C:\GO".
    ```

## 2. Installation

gin

    
    go get -u github.com/gin-gonic/gin
    

Excelize

    
    go get github.com/360EntSecGroup-Skylar/excelize
    

## 3. Run the program

    Open your Command Prompt, go to the project directory (For example, C:\GO\bin\src\Aiculus-Dashboard).
    
    
    Type in 
    go build
    
    Then type in 
    go run main.go DAO.go routes.go
    


# Full documentation

    Please refer to "Technical Report Aiculus.pdf". 
