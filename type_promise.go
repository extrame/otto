package otto

func (rt *runtime) newPromiseObject(value *promise) *object {
	obj := rt.newObject()
	obj.class = classPromiseName // TODO Should this be something else?
	obj.objectClass = classPromise
	obj.value = value

	return obj
}
