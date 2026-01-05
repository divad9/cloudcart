/**
 * Category Controller
 */
const Category = require('../models/Category');
const Product = require('../models/Product');

// Helper function to generate slug
const generateSlug = (name) => {
  return name
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)/g, '');
};

/**
 * Get all categories
 */
exports.getAllCategories = async (req, res) => {
  try {
    const { include_products = false } = req.query;

    const options = {
      where: { is_active: true },
      order: [['name', 'ASC']]
    };

    // Include subcategories
    options.include = [
      {
        model: Category,
        as: 'subcategories',
        where: { is_active: true },
        required: false
      }
    ];

    // Optionally include products
    if (include_products === 'true') {
      options.include.push({
        model: Product,
        as: 'products',
        where: { is_active: true },
        required: false,
        attributes: ['id', 'name', 'slug', 'price', 'images']
      });
    }

    const categories = await Category.findAll(options);

    res.json({ categories });

  } catch (error) {
    console.error('Error fetching categories:', error);
    res.status(500).json({ error: 'Failed to fetch categories' });
  }
};

/**
 * Get single category
 */
exports.getCategory = async (req, res) => {
  try {
    const { id } = req.params;

    const where = isNaN(id) ? { slug: id } : { id: parseInt(id) };

    const category = await Category.findOne({
      where: { ...where, is_active: true },
      include: [
        {
          model: Category,
          as: 'subcategories',
          where: { is_active: true },
          required: false
        },
        {
          model: Product,
          as: 'products',
          where: { is_active: true },
          required: false,
          limit: 20
        }
      ]
    });

    if (!category) {
      return res.status(404).json({ error: 'Category not found' });
    }

    res.json({ category });

  } catch (error) {
    console.error('Error fetching category:', error);
    res.status(500).json({ error: 'Failed to fetch category' });
  }
};

/**
 * Create category
 */
exports.createCategory = async (req, res) => {
  try {
    const {
      name,
      description,
      image_url,
      parent_id
    } = req.body;

    const slug = generateSlug(name);

    // Check if slug exists
    const existing = await Category.findOne({ where: { slug } });
    if (existing) {
      return res.status(400).json({ error: 'Category with this name already exists' });
    }

    const category = await Category.create({
      name,
      slug,
      description,
      image_url,
      parent_id
    });

    res.status(201).json({
      message: 'Category created successfully',
      category
    });

  } catch (error) {
    console.error('Error creating category:', error);
    res.status(500).json({ error: 'Failed to create category' });
  }
};

/**
 * Update category
 */
exports.updateCategory = async (req, res) => {
  try {
    const { id } = req.params;
    const updateData = req.body;

    const category = await Category.findByPk(id);
    if (!category) {
      return res.status(404).json({ error: 'Category not found' });
    }

    if (updateData.name && updateData.name !== category.name) {
      updateData.slug = generateSlug(updateData.name);
    }

    await category.update(updateData);

    res.json({
      message: 'Category updated successfully',
      category
    });

  } catch (error) {
    console.error('Error updating category:', error);
    res.status(500).json({ error: 'Failed to update category' });
  }
};

/**
 * Delete category
 */
exports.deleteCategory = async (req, res) => {
  try {
    const { id } = req.params;

    const category = await Category.findByPk(id);
    if (!category) {
      return res.status(404).json({ error: 'Category not found' });
    }

    await category.update({ is_active: false });

    res.json({
      message: 'Category deleted successfully'
    });

  } catch (error) {
    console.error('Error deleting category:', error);
    res.status(500).json({ error: 'Failed to delete category' });
  }
};