# Expenditure-Tracking-Bot

Telegram bot to submit your daily expenses and track your monthly expenses.

## Features

### Submit daily expenses
- Submit daily expenses to the bot in a chat style, by entering your expense details such as the name, date, amount and category of the transaction.
- Get your entered expense summarised in the response to confirm the expense you entered.

https://github.com/user-attachments/assets/c8ddb341-2ca2-4152-97c4-9a563640d4c7

### Pre-fill common expenses
- Set pre-filled expense so that you do not need to enter the same expenses often.
- Select the pre-filled expenses with saved settings such as the expense category.

https://github.com/user-attachments/assets/df75c3f9-4294-48a9-9c7d-980fc670f249

### View summary of monthly expense
- View the summary of your expense in the month, with the breakdown of your expense per category.

<img width="843" height="230" alt="image" src="https://github.com/user-attachments/assets/b7197b74-5cd7-4347-a32f-f195bb07bf10" />

## Running the program

To run the program
```
go run main
```

Make commands

```
make install     # set up virtual environment & install deps
make run         # run the model server
make freeze      # update requirements.txt
make clean       # delete virtual environment
make start       # install + run in one go
```
