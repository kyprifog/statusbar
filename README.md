
# Statusbar

![Screenshot](https://github.com/kyprifog/statusbar/blob/master/images/screenshot.png)

Statusbar is a tool for tracking your progress in periodic tasks.

# Usage

```
go build
cp .bars.example.yaml ~/.bars.yaml
./statusbar
```

Edit ~/.bars.yaml as needed.  Press Escape twice to exit

# Philosopy
The idea behind status bar is that you may not have a regularly scheduled time for these events or maybe you have a hectic life (like having kids).  This increments the bottom blue bar once every day, and the goal of status bar is to "race the clock" by making sure you increment the other bars appropriately.  The bars will gradually change color as you fall behind reminding you to attend to those items without having to regularly schedule them.

You can change the increment to change the frequency you want to do that item.  For example, to do it once a week, set increment to 7.  3 times increment to 2, etc.
