package domain

import (
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

// ============================================================================
// CUID HOOKS FOR AUTO-GENERATION
// ============================================================================

// BeforeCreate hook generates a CUID if ID is empty
// Attach this to any model using CUID as primary key

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = cuid.New()
	}
	return nil
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = cuid.New()
	}
	return nil
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = cuid.New()
	}
	return nil
}

func (c *Class) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = cuid.New()
	}
	return nil
}

func (cd *ClassDetail) BeforeCreate(tx *gorm.DB) error {
	if cd.ID == "" {
		cd.ID = cuid.New()
	}
	return nil
}

func (cp *ClassProgress) BeforeCreate(tx *gorm.DB) error {
	if cp.ID == "" {
		cp.ID = cuid.New()
	}
	return nil
}

func (cert *Certificate) BeforeCreate(tx *gorm.DB) error {
	if cert.ID == "" {
		cert.ID = cuid.New()
	}
	return nil
}

func (d *Discussion) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = cuid.New()
	}
	return nil
}

func (r *Reply) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = cuid.New()
	}
	return nil
}

func (s *Showcase) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = cuid.New()
	}
	return nil
}

func (sc *ShowcaseComment) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == "" {
		sc.ID = cuid.New()
	}
	return nil
}

func (cp *CommunityPoint) BeforeCreate(tx *gorm.DB) error {
	if cp.ID == "" {
		cp.ID = cuid.New()
	}
	return nil
}

func (ul *UserLevel) BeforeCreate(tx *gorm.DB) error {
	if ul.ID == "" {
		ul.ID = cuid.New()
	}
	return nil
}
