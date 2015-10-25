# go-concurrency

Simple reminder of component's responsibility

## A producer
  Which create order for thurty Xebian, who register the order into one redis DB with a limited TTL.

## A Bar
  Which must be called by the server in order to honour the order. It change the order into Redis with the information
  of whom performed it.
  
## A consumer
  which receive the beverages and if all is ok in redis, will increase the count of +1
  
## A Waiter
  receive order and call the Bar in order to serve the beverage.
