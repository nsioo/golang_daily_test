## Change Log
#### version 2.0.0 2019-10-28
##### New Features:
* Expose Sender. User could make Sender with custom config for special scenario.

##### Improvements
* Truncate the big log for avoiding data loss.
* Optimize metrics counter.
* GracefullyExit() just stop all Sender, which no longer make Send() return failure.
* Occurred error when write data to socket, the error detail could print to stderr.



