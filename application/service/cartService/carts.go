package cartService

import (
	"CipProject/application/dao/cartDao"
	"CipProject/application/domain"
)

// Add 加入购物车
func Add(cart domain.Cart) error {
	err := cartDao.Save(cart)
	return err
}

// GetCartByUsername 根据用户名查询购物车信息
func GetCartByUsername(username string) ([]domain.CartInfoDTO, error) {
	carts, err := cartDao.GetCartByUsername(username)
	return carts, err
}

// GetById 根据id查询购物车数据
func GetById(id int64) (domain.Cart, error) {
	cart, err := cartDao.GetById(id)
	return cart, err
}

// RemoveCart 删除购物车中指定记录
func RemoveCart(id int64) error {
	err := cartDao.Delete(id)
	return err
}
