function memoize(fn, options) {
  var cache = options && options.cache
    ? options.cache
    : cacheDefault;

  var serializer = options && options.serializer
    ? options.serializer
    : serializerDefault;

  var strategy = options && options.strategy
    ? options.strategy
    : strategyDefault;

  return strategy(fn, {
    cache: cache,
    serializer: serializer
  });
}

function isPrimitive(value) {
  return value == null || typeof value === 'number' || typeof value === 'boolean';
}

function monadic(fn, cache, serializer, arg) {
  var cacheKey = isPrimitive(arg) ? arg : serializer(arg);

  var computedValue = cache.get(cacheKey);
  if (typeof computedValue === 'undefined') {
    computedValue = fn.call(this, arg);
    cache.set(cacheKey, computedValue);
  }

  return computedValue;
}

function variadic(fn, cache, serializer) {
  var args = Array.prototype.slice.call(arguments, 3);
  var cacheKey = serializer(args);

  var computedValue = cache.get(cacheKey);
  if (typeof computedValue === 'undefined') {
    computedValue = fn.apply(this, args);
    cache.set(cacheKey, computedValue);
  }

  return computedValue;
}

function assemble(fn, context, strategy, cache, serialize) {
  return strategy.bind(
    context,
    fn,
    cache,
    serialize
  );
}

function strategyDefault(fn, options) {
  var strategy = fn.length === 1 ? monadic : variadic;

  return assemble(
    fn,
    this,
    strategy,
    options.cache.create(),
    options.serializer
  );
}

function strategyVariadic(fn, options) {
  var strategy = variadic;

  return assemble(
    fn,
    this,
    strategy,
    options.cache.create(),
    options.serializer
  );
}

function strategyMonadic(fn, options) {
  var strategy = monadic;

  return assemble(
    fn,
    this,
    strategy,
    options.cache.create(),
    options.serializer
  );
}

function serializerDefault() {
  return JSON.stringify(arguments);
}

function ObjectWithoutPrototypeCache() {
  this.cache = Object.create(null);
}

ObjectWithoutPrototypeCache.prototype.has = function(key) {
  return (key in this.cache);
};

ObjectWithoutPrototypeCache.prototype.get = function(key) {
  return this.cache[key];
};

ObjectWithoutPrototypeCache.prototype.set = function(key, value) {
  this.cache[key] = value;
};

var cacheDefault = {
  create: function create() {
    return new ObjectWithoutPrototypeCache();
  }
};

var strategies = {
  variadic: strategyVariadic,
  monadic: strategyMonadic
};

function getCurrentDateTime(includeDate, includeTime, includeSecond, includeMillisecond) {
    const now = new Date();
    let result = '';

    if (includeDate) {
        result += now.toLocaleDateString();
    }

    if (includeTime) {
        if (result !== '') result += ' ';
        result += now.toLocaleTimeString();
    }

    if (includeSecond) {
        if (result !== '') result += ' ';
        result += now.getSeconds();
    }

    if (includeMillisecond) {
        if (result !== '') result += ' ';
        result += now.getMilliseconds();
    }

    return result;
}

function measureExecutionTime(func) {
    const start = performance.now();
    func();
    const end = performance.now();
    return end - start;
}
