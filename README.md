# Debt Bot for Telegram [![Build Status](https://img.shields.io/travis/cloudedcat/debt-bot)](https://travis-ci.org/cloudedcat/debt-bot) [![Codecov](https://img.shields.io/codecov/c/github/cloudedcat/debt-bot)](https://codecov.io/gh/cloudedcat/debt-bot)

Telegram bot for recording restaurant bills and split them with participants efficiently

![Add bot](img/add_bot_example.png?raw=true) ![Reg command](img/reg_example.png?raw=true) ![Share command](img/share_example.png?raw=true)

## Commands

* `reg` - registers new participant.
* `list` - lists of participants
* `share <Amount> [in <Restaurant>] with <@username1> <@username2>...` - means that message sender paid \<Amount> in \<Restaurant> and want to split the bill equally between pointed participants and himself/herself. E.g.:
    > `@A` sent: `share 12 with @B @C` \
    > So, `@A` paid 12 for `@A`, `@B` and `@C` \
    > Since that `@B` and `@C` owe to `@A` 4

* `calc` - summarizes all debts and get optipmal way to get debts back. E.g.:
    > A: `/share 12 with @B @C` \
    > B: `/share 12 with @A @C` \
    > C: `/calc` \
    > Bot:
    >> `list of debts:` \
    >> `@C -> @A - 4.00` \
    >> `@C -> @B - 4.00`

* `history` - shows personal history of debts
