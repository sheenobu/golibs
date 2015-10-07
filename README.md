# sheenobu/golibs

Code that doesn't yet deserve it's own repo (yet)

## Packages

### log

extended logging functionality via log15 and golang.org/x/net/context:

	ctx := log.NewContext(ctx, params)

	log.Log(ctx).Debug(...)
	log.Log(ctx).Error(

	l := log.Log(ctx)
	l.Debug("X")

	l := log.Log(cxt).Extend(params)
	l.Debug("X")

### apps

Nested application and subprocess management


            pctx
            /
         app
         /
pctx  ctx
  \   / \
   app  process
     \
     ctx
    /    \
process  process



### dispatch

Simple pub/sub style dispatcher with named channels.

