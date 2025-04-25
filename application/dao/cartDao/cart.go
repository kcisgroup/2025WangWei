package cartDao

import (
	"CipProject/application/dao"
	"CipProject/application/domain"
)

// Save 新增数据
func Save(cart domain.Cart) error {
	db := dao.GetDB()
	if err := db.Create(&cart).Error; err != nil {
		return err
	}
	return nil
}

// GetCartByUsername 根据购物车拥有者的用户名查询购物车详情
func GetCartByUsername(username string) ([]domain.CartInfoDTO, error) {
	var carts []domain.CartInfoDTO
	db := dao.GetDB()

	// carts表左外连接code_records表
	err := db.Table("carts").Select("carts.*,code_records.filename,code_records.version,code_records.tag,"+
		"code_records.price,code_records.description").Joins("left join code_records on "+
		"carts.code_id = code_records.id").Where("owner = ?", username).Scan(&carts).Error
	return carts, err
}

func GetById(id int64) (domain.Cart, error) {
	var cart = domain.Cart{ID: id}
	db := dao.GetDB()
	err := db.First(&cart).Error
	return cart, err
}

// Delete 删除指定记录
func Delete(id int64) error {
	db := dao.GetDB()
	result := db.Delete(&domain.Cart{}, id)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
