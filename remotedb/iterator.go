package remotedb

import (
	tmdb "github.com/line/tm-db/v2"
	protodb "github.com/line/tm-db/v2/remotedb/proto"
)

func makeIterator(dic protodb.DB_IteratorClient) tmdb.Iterator {
	itr := &iterator{dic: dic}
	itr.Next() // We need to call Next to prime the iterator
	return itr
}

func makeReverseIterator(dric protodb.DB_ReverseIteratorClient) tmdb.Iterator {
	rItr := &reverseIterator{dric: dric}
	rItr.Next() // We need to call Next to prime the iterator
	return rItr
}

type reverseIterator struct {
	dric protodb.DB_ReverseIteratorClient
	cur  *protodb.Iterator
	err  error
}

var _ tmdb.Iterator = (*iterator)(nil)

// Valid implements Iterator.
func (rItr *reverseIterator) Valid() bool {
	return rItr.cur != nil && rItr.cur.Valid && rItr.err == nil
}

// Next implements Iterator.
func (rItr *reverseIterator) Next() {
	var err error
	rItr.cur, err = rItr.dric.Recv()
	if err != nil {
		rItr.err = err
	}
}

// Key implements Iterator.
func (rItr *reverseIterator) Key() []byte {
	rItr.assertIsValid()
	return rItr.cur.Key
}

// Value implements Iterator.
func (rItr *reverseIterator) Value() []byte {
	rItr.assertIsValid()
	return rItr.cur.Value
}

// Error implements Iterator.
func (rItr *reverseIterator) Error() error {
	return rItr.err
}

// Close implements Iterator.
func (rItr *reverseIterator) Close() error {
	return nil
}

func (rItr *reverseIterator) assertIsValid() {
	if !rItr.Valid() {
		panic("iterator is invalid")
	}
}

// iterator implements the tmdb.Iterator by retrieving
// streamed iterators from the remote backend as
// needed. It is NOT safe for concurrent usage,
// matching the behavior of other iterators.
type iterator struct {
	dic protodb.DB_IteratorClient
	cur *protodb.Iterator
	err error
}

var _ tmdb.Iterator = (*iterator)(nil)

// Valid implements Iterator.
func (itr *iterator) Valid() bool {
	return itr.cur != nil && itr.cur.Valid && itr.err == nil
}

// Next implements Iterator.
func (itr *iterator) Next() {
	var err error
	itr.cur, err = itr.dic.Recv()
	if err != nil {
		itr.err = err
	}
}

// Key implements Iterator.
func (itr *iterator) Key() []byte {
	itr.assertIsValid()
	return itr.cur.Key
}

// Value implements Iterator.
func (itr *iterator) Value() []byte {
	itr.assertIsValid()
	return itr.cur.Value
}

// Error implements Iterator.
func (itr *iterator) Error() error {
	return itr.err
}

// Close implements Iterator.
func (itr *iterator) Close() error {
	return itr.dic.CloseSend()
}

func (itr *iterator) assertIsValid() {
	if !itr.Valid() {
		panic("iterator is invalid")
	}
}
