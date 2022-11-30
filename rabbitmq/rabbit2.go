package rabbitmq

// type connectionWrapper struct {
// 	amqpAddr                string
// 	amqpConnection          *amqp.Connection
// 	amqpConnectionLock      sync.Mutex
// 	connected               chan struct{}
// 	publishBindingsLock     sync.RWMutex
// 	publishBindingsPrepared map[string]struct{}
// 	closing                 chan struct{}
// 	closed                  bool
// 	publishingMsg           sync.WaitGroup
// 	subscribingWg           sync.WaitGroup
// }

// func NewConnectionWrapper() (*connectionWrapper, error) {
// 	pubSub := &connectionWrapper{
// 		amqpAddr:  os.Getenv("RABBITMQ_ADDRESS"),
// 		closing:   make(chan struct{}),
// 		connected: make(chan struct{}),
// 	}

// 	if err := pubSub.connect(); err != nil {
// 		return nil, err
// 	}
// 	pubSub.handleConnectionClose()

// 	return pubSub, nil
// }

// func (c *connectionWrapper) handleConnectionClose() {
// 	for {
// 		log.Debug().Msg("handleConnectionClose is waiting for p.connected")
// 		<-c.connected
// 		log.Debug().Msg("handleConnectionClose is for connection or Pub/Sub close")

// 		notifyCloseConnection := c.amqpConnection.NotifyClose(make(chan *amqp.Error))

// 		select {
// 		case <-c.closing:
// 			log.Debug().Msg("Stopping handleConnectionClose")
// 			return
// 		case err := <-notifyCloseConnection:
// 			c.connected = make(chan struct{})
// 			log.Error().Msgf("Received close notification from AMQP, reconnecting %v", err)
// 			if err := c.reconnect(); err != nil {
// 				log.Debug().Msgf("reconnect failed %v", err)
// 				// c.Close()
// 			}
// 		}
// 	}
// }

// func (c *connectionWrapper) reconnect() error {
// 	b := &backoff.ExponentialBackOff{
// 		InitialInterval:     time.Duration(5) * time.Millisecond,
// 		RandomizationFactor: backoff.DefaultRandomizationFactor,
// 		Multiplier:          backoff.DefaultMultiplier,
// 		MaxInterval:         time.Duration(40) * time.Millisecond,
// 		MaxElapsedTime:      time.Duration(600) * time.Millisecond,
// 		Clock:               backoff.SystemClock,
// 	}

// 	if b.InitialInterval == 0 {
// 		b.InitialInterval = backoff.DefaultInitialInterval
// 	}

// 	if b.MaxInterval == 0 {
// 		b.MaxInterval = time.Minute
// 	}

// 	if b.MaxElapsedTime == 0 {
// 		b.MaxElapsedTime = 3 * time.Minute
// 	}
// 	now, attempt := time.Now(), 0
// 	b.Reset()
// 	err := backoff.Retry(func() (err error) {
// 		defer func() {
// 			log.Info().Msgf("AMQP connection open attempt %d cumulative opening time %v: %v",
// 				attempt, time.Since(now), err,
// 			)
// 		}()
// 		attempt++
// 		return c.connect()
// 	}, b)
// 	return err
// }

// func (c *connectionWrapper) connect() error {
// 	c.amqpConnectionLock.Lock()
// 	defer c.amqpConnectionLock.Unlock()
// 	var err error

// 	c.amqpConnection, err = amqp.Dial(c.amqpAddr)
// 	utils.FailOnError("rabbitmq", err)

// 	close(c.connected)

// 	utils.LogWithInfo("rabbitmq", "connected to rabbitMQ")

// 	return nil
// }

// type RabbitMQClient struct {
// 	channel *amqp.Channel
// }
